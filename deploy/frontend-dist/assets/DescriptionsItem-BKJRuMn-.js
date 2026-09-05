import{$r as e,Dr as t,Mi as n,Nn as r,Qr as i,Rr as a,Vi as o,Xn as s,Zn as c,Zr as l,cr as u,ei as d,ii as f,ni as p,ri as m,ti as h,zi as g}from"./router-BXCsM_PA.js";import{n as _}from"./Space-BC0LoIqy.js";import{s as v}from"./index-BCfmTlJD.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=l([i(`descriptions`,{fontSize:`var(--n-font-size)`},[i(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),i(`descriptions-table-wrapper`,[i(`descriptions-table`,[i(`descriptions-table-row`,[i(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),i(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),h(`bordered`,[i(`descriptions-table-wrapper`,[i(`descriptions-table`,[i(`descriptions-table-row`,[l(`&:last-child`,[i(`descriptions-table-content`,{paddingBottom:0})])])])])]),d(`left-label-placement`,[i(`descriptions-table-content`,[l(`> *`,{verticalAlign:`top`})])]),d(`left-label-align`,[l(`th`,{textAlign:`left`})]),d(`center-label-align`,[l(`th`,{textAlign:`center`})]),d(`right-label-align`,[l(`th`,{textAlign:`right`})]),d(`bordered`,[i(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[i(`descriptions-table`,[i(`descriptions-table-row`,[l(`&:not(:last-child)`,[i(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),i(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),i(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[l(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),i(`descriptions-table-content`,[l(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),i(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),i(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[i(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[i(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[i(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),i(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[e(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),e(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),i(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),m(i(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),f(i(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=g({name:`Descriptions`,props:Object.assign(Object.assign({},r.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:i,inlineThemeDisabled:a,mergedComponentPropsRef:o}=c(e),l=n(()=>e.size||o?.value?.Descriptions?.size||`medium`),u=r(`Descriptions`,`-descriptions`,b,v,e,i),d=n(()=>{let{bordered:t}=e,n=l.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:s,thTextColor:c,thFontWeight:d,tdTextColor:f,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[p(`fontSize`,n)]:S,[p(t?`thPaddingBordered`:`thPadding`,n)]:C,[p(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=u.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":d,"--n-line-height":x,"--n-th-text-color":c,"--n-td-text-color":f,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":s,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),f=a?s(`descriptions`,n(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=l.value[0],t}),d,e):void 0;return{mergedClsPrefix:i,cssVars:a?void 0:d,themeClass:f?.themeClass,onRender:f?.onRender,compitableColumn:t(e,[`columns`,`column`]),inlineThemeDisabled:a,mergedSize:l}},render(){let e=this.$slots.default,t=e?u(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:i,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:d,title:f,cssVars:p,mergedClsPrefix:m,separator:h,onRender:g}=this;g?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,a)=>{let c=t.props||{},l=v.length-1===a,u=[`label`in c?c.label:y(t,`label`)],f=[y(t)],p=c.span||1,g=e.span;e.span+=p;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(o(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:1,style:_},u),o(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:l?(i-g)*2+1:p*2-1,style:b},f)):e.row.push(o(`td`,{class:`${m}-descriptions-table-content`,colspan:l?(i-g)*2:p*2},o(`span`,{class:[`${m}-descriptions-table-content__label`,r],style:_},[...u,h&&o(`span`,{class:`${m}-descriptions-separator`},h)]),o(`span`,{class:[`${m}-descriptions-table-content__content`,n],style:b},f)));else{let t=l?(i-g)*2:p*2;e.row.push(o(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:t,style:_},u)),e.secondRow.push(o(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:t,style:b},f))}return(e.span>=i||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>o(`tr`,{class:`${m}-descriptions-table-row`},e));return o(`div`,{style:p,class:[`${m}-descriptions`,this.themeClass,`${m}-descriptions--${s}-label-placement`,`${m}-descriptions--${c}-label-align`,`${m}-descriptions--${l}-size`,d&&`${m}-descriptions--bordered`]},f||this.$slots.header?o(`div`,{class:`${m}-descriptions-header`},f||_(this,`header`)):null,o(`div`,{class:`${m}-descriptions-table-wrapper`},o(`table`,{class:`${m}-descriptions-table`},o(`tbody`,null,s===`top`&&o(`tr`,{class:`${m}-descriptions-table-row`,style:{visibility:`collapse`}},a(i*2,o(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=g({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};